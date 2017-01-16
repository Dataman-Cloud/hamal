export class releaseBackend {
  constructor($resource, BACKEND_URL_BASE) {
    'ngInject';

    this.$resource = $resource;
    this.BACKEND_URL_BASE = BACKEND_URL_BASE;
  }

  projects() {
    return this.$resource(`${this.BACKEND_URL_BASE.defaultBase}/v1/hamal/projects/:project`, {project: '@project'});
  }

  applications() {
    return this.$resource(`${this.BACKEND_URL_BASE.defaultBase}/v1/hamal/apps/:app`, {app: '@app'});
  }
}
