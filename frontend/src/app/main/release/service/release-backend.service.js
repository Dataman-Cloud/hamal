export class releaseBackend {
  constructor($resource, BACKEND_URL_BASE) {
    'ngInject';

    this.$resource = $resource;
    this.BACKEND_URL_BASE = BACKEND_URL_BASE;
  }

}
