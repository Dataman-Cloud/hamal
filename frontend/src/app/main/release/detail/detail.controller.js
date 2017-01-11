/**
 * Created by my9074 on 2017/1/11.
 */
export class DetailController {
  constructor(releaseBackend, $stateParams, $q) {
    'ngInject';

    this.releaseBackend = releaseBackend;
    this.$stateParams = $stateParams;
    this.$q = $q;
    this.project = {};
    this.apps = [];
    this.activate();
  }

  activate() {
    this.getProject()
  }

  getProject() {
    this.releaseBackend.projects().get({project: this.$stateParams.project}, data => {
      this.project = data.data;

      if (Array.isArray(this.project.applications)) {
        let prom = [];

        this.project.applications.forEach(app => {
          prom.push(this.getApplication(app.app_id))
        });

        this.$q.all(prom).then(data => {
          this.apps = data.map(item => item.data);
        })
      }
    })
  }

  getApplication(appId) {
    return this.releaseBackend.applications().get({app: appId}).$promise
  }
}
