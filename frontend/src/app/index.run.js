export function runBlock ($rootScope, $state, $stateParams) {
  'ngInject';

  $rootScope.$state = $state;
  $rootScope.$stateParams = $stateParams;
  $rootScope.keys = Object.keys;
}
