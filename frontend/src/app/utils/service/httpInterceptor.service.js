export function httpInterceptor($q, $injector, BACKEND_URL_BASE) {
  'ngInject';

  let notification = null;
  let getNotification = function () {
    if (!notification) {
      notification = $injector.get('Notification');
    }
    return notification;
  };

  return {
    // optional method
    'request': function (config) {
      // do something on success
      return config;
    },
    // optional method
    'requestError': function (rejection) {
      // do something on error
      return $q.reject(rejection);
    },

    // optional method
    'response': function (response) {
      // do something on success
      return response;
    },

    // optional method
    'responseError': function (rejection) {
      // do something on error
      let msg = '';
      if (-1 === rejection.status) {
        msg = `连接后端服务器异常 </br>
               请确认配置 ${BACKEND_URL_BASE.defaultBase}`;
      } else if (rejection.data) {
        msg = rejection.data;
      } else {
        msg = rejection.statusText;
      }
      getNotification().error(msg);
      return $q.reject(rejection);
    }

  }
}
