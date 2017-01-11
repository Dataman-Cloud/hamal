export function config ($logProvider, $locationProvider, $interpolateProvider, $httpProvider, cfpLoadingBarProvider) {
  'ngInject';
  // Enable log
  $logProvider.debugEnabled(true);

  // Set options third-party lib
  $locationProvider.html5Mode(true);

  $interpolateProvider.startSymbol('{/');
  $interpolateProvider.endSymbol('/}');
  $httpProvider.interceptors.push('httpInterceptor');

  cfpLoadingBarProvider.includeSpinner = false;
}
