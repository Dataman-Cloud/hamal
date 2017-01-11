/**
 * Created by my9074 on 2017/1/11.
 */
export class ListController {
  constructor(releaseBackend) {
    'ngInject';

    this.releaseBackend = releaseBackend;
    this.projects = [];
    this.activate();
  }

  activate() {
    this.listProject();
  }

  listProject() {
    this.releaseBackend.projects().get(data => {
      this.projects = data.data;
    })
  }
}
