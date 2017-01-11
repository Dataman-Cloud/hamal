export function routerConfig($stateProvider, $urlRouterProvider) {
  'ngInject';
  $stateProvider
    .state('home', {
      templateUrl: 'app/main/main.html',
      controller: 'MainController',
      controllerAs: 'main'
    })
    .state('home.list', {
      url: '/list',
      templateUrl: 'app/main/release/list/list.html',
      controller: 'ListController',
      controllerAs: 'vm'
    })
    .state('home.detail', {
    url: '/detail',
    templateUrl: 'app/main/release/detail/detail.html',
    controller: 'DetailController',
    controllerAs: 'vm'
  });

  $urlRouterProvider.otherwise($injector => {
    let $state = $injector.get('$state');
    $state.go('home.list');
  });
}
