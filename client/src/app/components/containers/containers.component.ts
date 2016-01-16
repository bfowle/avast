import {Component, OnInit, OnDestroy} from 'angular2/core';
import {Router} from 'angular2/router';
import {AppStore} from 'angular2-redux';

import {DockerActions} from '../../actions/docker.actions';
import {DockerContainer} from './interfaces/docker-container';

@Component({
  selector: 'avast-containers',
  template: require('./containers.component.html'),
  styles: [
    require('./containers.component.css')
  ]
})

export class ContainersComponent implements OnInit {
  public containers: DockerContainer[];
  public selectedContainer: DockerContainer;
  private isFetchingContainers: boolean = false;
  private unsubscribe: Function;

  constructor(private _router: Router,
              private _appStore: AppStore,
              private _dockerActions: DockerActions) {
  }

  ngOnInit() {
    this.unsubscribe = this._appStore.subscribe((state) => {
      this.containers = state.docker.containers;
      this.isFetchingContainers = state.docker.isFetchingContainers;
    });

    this._appStore.dispatch(this._dockerActions.fetchContainers());
  }

  onSelect(container: DockerContainer) {
    this.selectedContainer = container;
  }

  private ngOnDestroy() {
    this.unsubscribe();
  }
}
