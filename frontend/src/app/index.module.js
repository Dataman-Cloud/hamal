/* global malarkey:false, moment:false */

import { config } from './index.config';
import { routerConfig } from './index.route';
import { runBlock } from './index.run';
import { MainController } from './main/main.controller';
import { ListController } from './main/release/list/list.controller';
import { DetailController } from './main/release/detail/detail.controller';
import { httpInterceptor } from './utils/service/httpInterceptor.service';
import { releaseBackend } from './main/release/service/release-backend.service';
import { sizeFilter } from './filter/size.filter'

angular.module('frontend', ['ngAnimate', 'ngCookies', 'ngSanitize', 'ngMessages', 'ngAria', 'ngResource', 'ui.router', 'ngMaterial',
  'ui-notification', 'md.data.table', 'angular-loading-bar', 'diff-match-patch'])
  .constant('moment', moment)
  .constant('_', _)
  .constant('BACKEND_URL_BASE', {
    defaultBase: "http://192.168.1.24:5099"
  })
  .config(config)
  .config(routerConfig)
  .run(runBlock)
  .service('httpInterceptor', httpInterceptor)
  .service('releaseBackend', releaseBackend)
  .filter('size', sizeFilter)
  .controller('MainController', MainController)
  .controller('ListController', ListController)
  .controller('DetailController', DetailController);
