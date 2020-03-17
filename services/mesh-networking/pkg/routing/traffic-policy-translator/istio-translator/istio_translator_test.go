package istio_translator_test

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	core_types "github.com/solo-io/mesh-projects/pkg/api/core.zephyr.solo.io/v1alpha1/types"
	discovery_v1alpha1 "github.com/solo-io/mesh-projects/pkg/api/discovery.zephyr.solo.io/v1alpha1"
	discovery_types "github.com/solo-io/mesh-projects/pkg/api/discovery.zephyr.solo.io/v1alpha1/types"
	"github.com/solo-io/mesh-projects/pkg/api/networking.zephyr.solo.io/v1alpha1"
	networking_types "github.com/solo-io/mesh-projects/pkg/api/networking.zephyr.solo.io/v1alpha1/types"
	istio_networking "github.com/solo-io/mesh-projects/pkg/clients/istio/networking"
	mock_istio_networking "github.com/solo-io/mesh-projects/pkg/clients/istio/networking/mock"
	mock_core "github.com/solo-io/mesh-projects/pkg/clients/zephyr/discovery/mocks"
	"github.com/solo-io/mesh-projects/services/common"
	mock_mc_manager "github.com/solo-io/mesh-projects/services/common/multicluster/manager/mocks"
	istio_translator "github.com/solo-io/mesh-projects/services/mesh-networking/pkg/routing/traffic-policy-translator/istio-translator"
	mock_preprocess "github.com/solo-io/mesh-projects/services/mesh-networking/pkg/routing/traffic-policy-translator/preprocess/mocks"
	api_v1alpha3 "istio.io/api/networking/v1alpha3"
	client_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type testContext struct {
	clusterName            string
	meshObjKey             client.ObjectKey
	meshServiceObjKey      client.ObjectKey
	kubeServiceObjKey      client.ObjectKey
	mesh                   *discovery_v1alpha1.Mesh
	meshService            *discovery_v1alpha1.MeshService
	trafficPolicy          []*v1alpha1.TrafficPolicy
	computedVirtualService *client_v1alpha3.VirtualService
}

var _ = Describe("IstioTranslator", func() {
	var (
		ctrl                         *gomock.Controller
		istioTrafficPolicyTranslator istio_translator.IstioTranslator
		ctx                          context.Context
		mockDynamicClientGetter      *mock_mc_manager.MockDynamicClientGetter
		mockMeshClient               *mock_core.MockMeshClient
		mockMeshServiceClient        *mock_core.MockMeshServiceClient
		mockVirtualServiceClient     *mock_istio_networking.MockVirtualServiceClient
		mockDestinationRuleClient    *mock_istio_networking.MockDestinationRuleClient
		mockMeshServiceSelector      *mock_preprocess.MockMeshServiceSelector
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.TODO()
		mockDynamicClientGetter = mock_mc_manager.NewMockDynamicClientGetter(ctrl)
		mockMeshClient = mock_core.NewMockMeshClient(ctrl)
		mockMeshServiceClient = mock_core.NewMockMeshServiceClient(ctrl)
		mockVirtualServiceClient = mock_istio_networking.NewMockVirtualServiceClient(ctrl)
		mockMeshServiceSelector = mock_preprocess.NewMockMeshServiceSelector(ctrl)
		mockDestinationRuleClient = mock_istio_networking.NewMockDestinationRuleClient(ctrl)
		istioTrafficPolicyTranslator = istio_translator.NewIstioTrafficPolicyTranslator(
			mockDynamicClientGetter,
			mockMeshClient,
			mockMeshServiceClient,
			mockMeshServiceSelector,
			func(client client.Client) istio_networking.VirtualServiceClient {
				return mockVirtualServiceClient
			},
			func(client client.Client) istio_networking.DestinationRuleClient {
				return mockDestinationRuleClient
			},
		)
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("should translate TrafficPolicies into VirtualService and DestinationRule and upsert", func() {
		setupTestContext := func() *testContext {
			clusterName := "clusterName"
			meshObjKey := client.ObjectKey{Name: "mesh-name", Namespace: "mesh-namespace"}
			meshServiceObjKey := client.ObjectKey{Name: "mesh-service-name", Namespace: "mesh-service-namespace"}
			kubeServiceObjKey := client.ObjectKey{Name: "kube-service-name", Namespace: "kube-service-namespace"}
			meshServiceFederationMCDnsName := "multiclusterDNSname"
			meshService := &discovery_v1alpha1.MeshService{
				ObjectMeta: v1.ObjectMeta{
					Name:        meshServiceObjKey.Name,
					Namespace:   meshServiceObjKey.Namespace,
					ClusterName: clusterName,
				},
				Spec: discovery_types.MeshServiceSpec{
					Mesh: &core_types.ResourceRef{
						Name:      meshObjKey.Name,
						Namespace: meshObjKey.Namespace,
					},
					KubeService: &discovery_types.KubeService{
						Ref: &core_types.ResourceRef{
							Name:      kubeServiceObjKey.Name,
							Namespace: kubeServiceObjKey.Namespace,
							Cluster:   &types.StringValue{Value: clusterName},
						},
					},
					Federation: &discovery_types.Federation{
						MulticlusterDnsName: meshServiceFederationMCDnsName,
					},
				},
			}
			mesh := &discovery_v1alpha1.Mesh{
				Spec: discovery_types.MeshSpec{
					Cluster: &core_types.ResourceRef{
						Name: clusterName,
					},
					MeshType: &discovery_types.MeshSpec_Istio{
						Istio: &discovery_types.IstioMesh{},
					},
				},
			}
			trafficPolicy := []*v1alpha1.TrafficPolicy{{
				Spec: networking_types.TrafficPolicySpec{
					HttpRequestMatchers: []*networking_types.HttpMatcher{{}, {}, {}},
				}},
			}
			computedVirtualService := &client_v1alpha3.VirtualService{
				ObjectMeta: v1.ObjectMeta{
					Name:      meshService.Spec.GetKubeService().GetRef().GetName(),
					Namespace: meshService.Spec.GetKubeService().GetRef().GetNamespace(),
				},
				Spec: api_v1alpha3.VirtualService{
					Hosts: []string{meshServiceObjKey.Name},
					Http: []*api_v1alpha3.HTTPRoute{
						{
							Match: []*api_v1alpha3.HTTPMatchRequest{{}},
						},
						{
							Match: []*api_v1alpha3.HTTPMatchRequest{{}},
						},
						{
							Match: []*api_v1alpha3.HTTPMatchRequest{{}},
						},
					},
				},
			}
			mockMeshClient.EXPECT().Get(ctx, meshObjKey).Return(mesh, nil)
			mockDynamicClientGetter.EXPECT().GetClientForCluster(clusterName).Return(nil, true)
			// computed DestinationRule
			computedDestinationRule := &client_v1alpha3.DestinationRule{
				ObjectMeta: v1.ObjectMeta{
					Name:      meshService.Spec.GetKubeService().GetRef().GetName(),
					Namespace: meshService.Spec.GetKubeService().GetRef().GetNamespace(),
				},
				Spec: api_v1alpha3.DestinationRule{
					Host: meshServiceObjKey.Name,
					TrafficPolicy: &api_v1alpha3.TrafficPolicy{
						Tls: &api_v1alpha3.TLSSettings{
							Mode: api_v1alpha3.TLSSettings_ISTIO_MUTUAL,
						},
					},
				},
			}
			mockDestinationRuleClient.EXPECT().Upsert(ctx, computedDestinationRule).Return(nil)
			return &testContext{
				clusterName:            clusterName,
				meshObjKey:             meshObjKey,
				meshServiceObjKey:      meshServiceObjKey,
				kubeServiceObjKey:      kubeServiceObjKey,
				mesh:                   mesh,
				meshService:            meshService,
				trafficPolicy:          trafficPolicy,
				computedVirtualService: computedVirtualService,
			}
		}

		It("should upsert VirtualService", func() {
			testContext := setupTestContext()
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx,
				testContext.meshService,
				testContext.mesh,
				testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate CorsPolicy", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.CorsPolicy = &networking_types.CorsPolicy{
				AllowOrigins: []*networking_types.StringMatch{
					{MatchType: &networking_types.StringMatch_Exact{Exact: "exact"}},
					{MatchType: &networking_types.StringMatch_Prefix{Prefix: "prefix"}},
					{MatchType: &networking_types.StringMatch_Regex{Regex: "regex"}},
				},
				AllowMethods:     []string{"GET", "POST"},
				AllowHeaders:     []string{"Header1", "Header2"},
				ExposeHeaders:    []string{"ExposedHeader1", "ExposedHeader2"},
				MaxAge:           &types.Duration{Seconds: 1},
				AllowCredentials: &types.BoolValue{Value: false},
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.CorsPolicy = &api_v1alpha3.CorsPolicy{
					AllowOrigins: []*api_v1alpha3.StringMatch{
						{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "exact"}},
						{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "prefix"}},
						{MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "regex"}},
					},
					AllowMethods:     []string{"GET", "POST"},
					AllowHeaders:     []string{"Header1", "Header2"},
					ExposeHeaders:    []string{"ExposedHeader1", "ExposedHeader2"},
					MaxAge:           &types.Duration{Seconds: 1},
					AllowCredentials: &types.BoolValue{Value: false},
				}
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HeaderManipulation", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HeaderManipulation = &networking_types.HeaderManipulation{
				AppendRequestHeaders:  map[string]string{"a": "b"},
				RemoveRequestHeaders:  []string{"3", "4"},
				AppendResponseHeaders: map[string]string{"foo": "bar"},
				RemoveResponseHeaders: []string{"1", "2"},
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Headers = &api_v1alpha3.Headers{
					Request: &api_v1alpha3.Headers_HeaderOperations{
						Add:    map[string]string{"a": "b"},
						Remove: []string{"3", "4"},
					},
					Response: &api_v1alpha3.Headers_HeaderOperations{
						Add:    map[string]string{"foo": "bar"},
						Remove: []string{"1", "2"},
					},
				}
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate Mirror destination on same cluster", func() {
			testContext := setupTestContext()
			destName := "name"
			destNamespace := "namespace"
			destCluster := &types.StringValue{Value: testContext.clusterName}
			testContext.trafficPolicy[0].Spec.Mirror = &networking_types.Mirror{
				Destination: &core_types.ResourceRef{
					Name:      destName,
					Namespace: destNamespace,
					Cluster:   destCluster,
				},
				Percentage: 50,
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Mirror = &api_v1alpha3.Destination{
					Host: destName + "." + destNamespace,
				}
				httpRoute.MirrorPercentage = &api_v1alpha3.Percent{Value: 50.0}
			}
			backingMeshService := &discovery_v1alpha1.MeshService{
				Spec: discovery_types.MeshServiceSpec{
					KubeService: &discovery_types.KubeService{
						Ref: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
						},
					},
				},
			}
			mockMeshServiceSelector.
				EXPECT().
				GetBackingMeshService(ctx, destName, destNamespace, destCluster.GetValue()).
				Return(backingMeshService, nil)
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate Mirror destination on same *local* cluster", func() {
			testContext := setupTestContext()
			destName := "name"
			destNamespace := "namespace"
			testContext.meshService.Spec.GetKubeService().GetRef().GetCluster().Value = common.LocalClusterName
			testContext.trafficPolicy[0].Spec.Mirror = &networking_types.Mirror{
				Destination: &core_types.ResourceRef{
					Name:      destName,
					Namespace: destNamespace,
					// omit cluster to specify local cluster
				},
				Percentage: 50,
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Mirror = &api_v1alpha3.Destination{
					Host: destName + "." + destNamespace,
				}
				httpRoute.MirrorPercentage = &api_v1alpha3.Percent{Value: 50.0}
			}
			backingMeshService := &discovery_v1alpha1.MeshService{
				Spec: discovery_types.MeshServiceSpec{
					KubeService: &discovery_types.KubeService{
						Ref: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
						},
					},
				},
			}
			mockMeshServiceSelector.
				EXPECT().
				GetBackingMeshService(ctx, destName, destNamespace, "").
				Return(backingMeshService, nil)
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate Mirror destination on remote cluster", func() {
			testContext := setupTestContext()
			multiClusterDnsName := "multicluster-dns-name"
			destName := "name"
			destNamespace := "namespace"
			remoteClusterName := "remote-cluster"
			testContext.trafficPolicy[0].Spec.Mirror = &networking_types.Mirror{
				Destination: &core_types.ResourceRef{
					Name:      destName,
					Namespace: destNamespace,
					Cluster:   &types.StringValue{Value: remoteClusterName},
				},
				Percentage: 50,
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Mirror = &api_v1alpha3.Destination{
					Host: multiClusterDnsName,
				}
				httpRoute.MirrorPercentage = &api_v1alpha3.Percent{Value: 50.0}
			}
			backingMeshService := &discovery_v1alpha1.MeshService{
				Spec: discovery_types.MeshServiceSpec{
					KubeService: &discovery_types.KubeService{
						Ref: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
						},
					},
					Federation: &discovery_types.Federation{MulticlusterDnsName: multiClusterDnsName},
				},
			}
			mockMeshServiceSelector.
				EXPECT().
				GetBackingMeshService(ctx, destName, destNamespace, remoteClusterName).
				Return(backingMeshService, nil)
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate FaultInjection of type Abort", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.FaultInjection = &networking_types.FaultInjection{
				FaultInjectionType: &networking_types.FaultInjection_Abort_{
					Abort: &networking_types.FaultInjection_Abort{
						ErrorType: &networking_types.FaultInjection_Abort_HttpStatus{HttpStatus: 404},
					},
				},
				Percentage: 50,
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Fault = &api_v1alpha3.HTTPFaultInjection{
					Abort: &api_v1alpha3.HTTPFaultInjection_Abort{
						ErrorType:  &api_v1alpha3.HTTPFaultInjection_Abort_HttpStatus{HttpStatus: 404},
						Percentage: &api_v1alpha3.Percent{Value: 50},
					},
				}
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate FaultInjection of type Delay of type Fixed", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.FaultInjection = &networking_types.FaultInjection{
				FaultInjectionType: &networking_types.FaultInjection_Delay_{
					Delay: &networking_types.FaultInjection_Delay{
						HttpDelayType: &networking_types.FaultInjection_Delay_FixedDelay{
							FixedDelay: &types.Duration{Seconds: 2},
						},
					},
				},
				Percentage: 50,
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Fault = &api_v1alpha3.HTTPFaultInjection{
					Delay: &api_v1alpha3.HTTPFaultInjection_Delay{
						HttpDelayType: &api_v1alpha3.HTTPFaultInjection_Delay_FixedDelay{FixedDelay: &types.Duration{Seconds: 2}},
						Percentage:    &api_v1alpha3.Percent{Value: 50},
					},
				}
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate FaultInjection of type Delay of type Exponential", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.FaultInjection = &networking_types.FaultInjection{
				FaultInjectionType: &networking_types.FaultInjection_Delay_{
					Delay: &networking_types.FaultInjection_Delay{
						HttpDelayType: &networking_types.FaultInjection_Delay_ExponentialDelay{
							ExponentialDelay: &types.Duration{Seconds: 2},
						},
					},
				},
				Percentage: 50,
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Fault = &api_v1alpha3.HTTPFaultInjection{
					Delay: &api_v1alpha3.HTTPFaultInjection_Delay{
						HttpDelayType: &api_v1alpha3.HTTPFaultInjection_Delay_ExponentialDelay{ExponentialDelay: &types.Duration{Seconds: 2}},
						Percentage:    &api_v1alpha3.Percent{Value: 50},
					},
				}
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate Retries", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.Retries = &networking_types.RetryPolicy{
				Attempts:      5,
				PerTryTimeout: &types.Duration{Seconds: 2},
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Retries = &api_v1alpha3.HTTPRetry{
					Attempts:      5,
					PerTryTimeout: &types.Duration{Seconds: 2},
				}
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HeaderMatchers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				Headers: []*networking_types.HeaderMatcher{
					{
						Name:        "name1",
						Value:       "value1",
						Regex:       false,
						InvertMatch: false,
					},
					{
						Name:        "name2",
						Value:       "*",
						Regex:       true,
						InvertMatch: false,
					},
					{
						Name:        "name3",
						Value:       "[a-z]+",
						Regex:       true,
						InvertMatch: true,
					},
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					Headers: map[string]*api_v1alpha3.StringMatch{
						"name1": {MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "value1"}},
						"name2": {MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "*"}},
					},
					WithoutHeaders: map[string]*api_v1alpha3.StringMatch{
						"name3": {MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "[a-z]+"}},
					},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HttpMatcher exact path specifiers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				PathSpecifier: &networking_types.HttpMatcher_Regex{
					Regex: "*",
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "*"}},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HttpMatcher prefix path specifiers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				PathSpecifier: &networking_types.HttpMatcher_Prefix{
					Prefix: "prefix",
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "prefix"}},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate QueryParamMatchers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				QueryParameters: []*networking_types.QueryParameterMatcher{
					{
						Name:  "qp1",
						Value: "qpv1",
						Regex: false,
					},
					{
						Name:  "qp2",
						Value: "qpv2",
						Regex: true,
					},
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					QueryParams: map[string]*api_v1alpha3.StringMatch{
						"qp1": {
							MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "qpv1"},
						},
						"qp2": {
							MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "qpv2"},
						},
					},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HttpMatcher regex path specifiers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				PathSpecifier: &networking_types.HttpMatcher_Regex{
					Regex: "*",
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "*"}},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HttpMatcher prefix path specifiers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				PathSpecifier: &networking_types.HttpMatcher_Prefix{
					Prefix: "prefix",
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "prefix"}},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate HttpMatcher exact path specifiers", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers[0] = &networking_types.HttpMatcher{
				Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				PathSpecifier: &networking_types.HttpMatcher_Exact{
					Exact: "path",
				},
			}
			testContext.computedVirtualService.Spec.Http[0].Match = []*api_v1alpha3.HTTPMatchRequest{
				{
					Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: networking_types.HttpMethodValue_GET.String()}},
					Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "path"}},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate TrafficShift without subsets", func() {
			testContext := setupTestContext()
			destName := "name"
			destNamespace := "namespace"
			multiClusterDnsName := "multicluster-dns-name"
			destCluster := &types.StringValue{Value: "remote-cluster-1"}
			testContext.trafficPolicy[0].Spec.TrafficShift = &networking_types.MultiDestination{
				Destinations: []*networking_types.MultiDestination_WeightedDestination{
					{
						Destination: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
							Cluster:   destCluster,
						},
						Weight: 50,
					},
				},
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Route = []*api_v1alpha3.HTTPRouteDestination{
					{
						Destination: &api_v1alpha3.Destination{
							Host: multiClusterDnsName,
						},
						Weight: 50,
					},
				}
			}
			backingMeshService := &discovery_v1alpha1.MeshService{
				Spec: discovery_types.MeshServiceSpec{
					KubeService: &discovery_types.KubeService{
						Ref: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
						},
					},
					Federation: &discovery_types.Federation{MulticlusterDnsName: multiClusterDnsName},
				},
			}
			mockMeshServiceSelector.
				EXPECT().
				GetBackingMeshService(ctx, destName, destNamespace, destCluster.GetValue()).
				Return(backingMeshService, nil)
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should translate TrafficShift with subsets", func() {
			testContext := setupTestContext()
			destName := "name"
			destNamespace := "namespace"
			multiClusterDnsName := "multicluster-dns-name"
			destCluster := &types.StringValue{Value: "remote-cluster-1"}
			declaredSubset := map[string]string{"env": "dev", "version": "v1"}
			expectedSubsetName := "env:dev,version:v1,"
			testContext.trafficPolicy[0].Spec.TrafficShift = &networking_types.MultiDestination{
				Destinations: []*networking_types.MultiDestination_WeightedDestination{
					{
						Destination: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
							Cluster:   destCluster,
						},
						Subset: declaredSubset,
						Weight: 50,
					},
				},
			}
			for _, httpRoute := range testContext.computedVirtualService.Spec.Http {
				httpRoute.Route = []*api_v1alpha3.HTTPRouteDestination{
					{
						Destination: &api_v1alpha3.Destination{
							Host:   multiClusterDnsName,
							Subset: expectedSubsetName,
						},
						Weight: 50,
					},
				}
			}
			backingMeshService := &discovery_v1alpha1.MeshService{
				Spec: discovery_types.MeshServiceSpec{
					KubeService: &discovery_types.KubeService{
						Ref: &core_types.ResourceRef{
							Name:      destName,
							Namespace: destNamespace,
						},
					},
					Federation: &discovery_types.Federation{MulticlusterDnsName: multiClusterDnsName},
				},
			}
			existingDestRule := &client_v1alpha3.DestinationRule{}
			computedDestRule := &client_v1alpha3.DestinationRule{
				Spec: api_v1alpha3.DestinationRule{
					Subsets: []*api_v1alpha3.Subset{
						{
							Name:   expectedSubsetName,
							Labels: declaredSubset,
						},
					},
				},
			}
			mockMeshServiceSelector.
				EXPECT().
				GetBackingMeshService(ctx, destName, destNamespace, destCluster.GetValue()).
				Return(backingMeshService, nil)

			mockDynamicClientGetter.
				EXPECT().
				GetClientForCluster(destCluster.GetValue()).
				Return(nil, true)
			mockDestinationRuleClient.
				EXPECT().
				Get(ctx, client.ObjectKey{Name: destName, Namespace: destNamespace}).
				Return(existingDestRule, nil)
			mockDestinationRuleClient.
				EXPECT().
				Update(ctx, computedDestRule).
				Return(nil)
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should return error if multiple MeshServices found for name/namespace/cluster", func() {
			testContext := setupTestContext()
			destName := "name"
			destNamespace := "namespace"
			remoteClusterName := "remote-cluster"
			testContext.trafficPolicy[0].Spec.Mirror = &networking_types.Mirror{
				Destination: &core_types.ResourceRef{
					Name:      destName,
					Namespace: destNamespace,
					Cluster:   &types.StringValue{Value: remoteClusterName},
				},
				Percentage: 50,
			}
			err := eris.New("mesh-service-selector-error")
			mockMeshServiceSelector.
				EXPECT().
				GetBackingMeshService(ctx, destName, destNamespace, remoteClusterName).
				Return(nil, err)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError.ErrorMessage).To(ContainSubstring(err.Error()))
		})

		It("should translate HTTP RequestMatchers and order the resulting HTTPRoutes", func() {
			testContext := setupTestContext()
			labels := map[string]string{"env": "dev"}
			namespaces := []string{"n1", "n2"}
			testContext.trafficPolicy[0].Spec.SourceSelector = &core_types.Selector{
				Labels:     labels,
				Namespaces: namespaces,
			}
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers = []*networking_types.HttpMatcher{
				{
					PathSpecifier: &networking_types.HttpMatcher_Exact{
						Exact: "path",
					},
					Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_GET},
				},
				{
					Headers: []*networking_types.HeaderMatcher{
						{
							Name:        "name3",
							Value:       "[a-z]+",
							Regex:       true,
							InvertMatch: true,
						},
					},
					Method: &networking_types.HttpMethod{Method: networking_types.HttpMethodValue_POST},
				},
			}
			testContext.computedVirtualService.Spec.Http = []*api_v1alpha3.HTTPRoute{
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Method:          &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "GET"}},
							Uri:             &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "path"}},
							SourceLabels:    labels,
							SourceNamespace: namespaces[1],
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "POST"}},
							WithoutHeaders: map[string]*api_v1alpha3.StringMatch{
								"name3": {MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "[a-z]+"}},
							},
							SourceLabels:    labels,
							SourceNamespace: namespaces[1],
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Method:          &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "GET"}},
							Uri:             &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "path"}},
							SourceLabels:    labels,
							SourceNamespace: namespaces[0],
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "POST"}},
							WithoutHeaders: map[string]*api_v1alpha3.StringMatch{
								"name3": {MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "[a-z]+"}},
							},
							SourceLabels:    labels,
							SourceNamespace: namespaces[0],
						},
					},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should deterministically order HTTPRoutes according to decreasing specificity", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers = []*networking_types.HttpMatcher{
				{
					PathSpecifier: &networking_types.HttpMatcher_Exact{
						Exact: "exact-path",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Prefix{
						Prefix: "/prefix",
					},
					Method: &networking_types.HttpMethod{
						Method: networking_types.HttpMethodValue_GET,
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Exact{
						Exact: "exact-path",
					},
					Method: &networking_types.HttpMethod{
						Method: networking_types.HttpMethodValue_GET,
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Exact{
						Exact: "exact-path",
					},
					Method: &networking_types.HttpMethod{
						Method: networking_types.HttpMethodValue_PUT,
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Regex{
						Regex: "www*",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Prefix{
						Prefix: "/",
					},
					Headers: []*networking_types.HeaderMatcher{
						{
							Name:        "set-cookie",
							Value:       "foo=bar",
							InvertMatch: true,
						},
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Prefix{
						Prefix: "/",
					},
					Headers: []*networking_types.HeaderMatcher{
						{
							Name:        "content-type",
							Value:       "text/html",
							InvertMatch: false,
						},
					},
				},
			}
			testContext.computedVirtualService.Spec.Http = []*api_v1alpha3.HTTPRoute{
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Headers: map[string]*api_v1alpha3.StringMatch{
								"content-type": {MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "text/html"}},
							},
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "/"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "exact-path"}},
							Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "PUT"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "exact-path"}},
							Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "GET"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "exact-path"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "www*"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri:    &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "/prefix"}},
							Method: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "GET"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							WithoutHeaders: map[string]*api_v1alpha3.StringMatch{
								"set-cookie": {MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "foo=bar"}},
							},
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "/"}},
						},
					},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})

		It("should order longer prefixes, regexes, and exact URI matchers before shorter ones", func() {
			testContext := setupTestContext()
			testContext.trafficPolicy[0].Spec.HttpRequestMatchers = []*networking_types.HttpMatcher{
				{
					PathSpecifier: &networking_types.HttpMatcher_Exact{
						Exact: "short",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Exact{
						Exact: "longer",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Prefix{
						Prefix: "/short",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Prefix{
						Prefix: "/longer",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Regex{
						Regex: "short*",
					},
				},
				{
					PathSpecifier: &networking_types.HttpMatcher_Regex{
						Regex: "longer*",
					},
				},
			}
			testContext.computedVirtualService.Spec.Http = []*api_v1alpha3.HTTPRoute{
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "longer"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Exact{Exact: "short"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "longer*"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Regex{Regex: "short*"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "/longer"}},
						},
					},
				},
				{
					Match: []*api_v1alpha3.HTTPMatchRequest{
						{
							Uri: &api_v1alpha3.StringMatch{MatchType: &api_v1alpha3.StringMatch_Prefix{Prefix: "/short"}},
						},
					},
				},
			}
			mockVirtualServiceClient.
				EXPECT().
				Upsert(ctx, testContext.computedVirtualService).
				Return(nil)
			translatorError := istioTrafficPolicyTranslator.TranslateTrafficPolicy(
				ctx, testContext.meshService, testContext.mesh, testContext.trafficPolicy)
			Expect(translatorError).To(BeNil())
		})
	})
})