import {Component, OnInit, OnDestroy} from 'angular2/core';
import {RouteParams} from 'angular2/router';
import {AppStore} from 'angular2-redux';

import {ConsulActions} from '../../actions/consul.actions';
import {SwarmNode} from './interfaces/swarm-node';

@Component({
  selector: 'avast-node-detail',
  template: require('./node-detail.component.html'),
  styles: [
    require('./node-detail.component.css')
  ]
})

export class NodeDetailComponent implements OnInit {
  public node: SwarmNode;
  private isFetchingNode: boolean = false;
  private unsubscribe: Function;

  constructor(private _routeParams: RouteParams,
              private _appStore: AppStore,
              private _consulActions: ConsulActions) {
  }

  ngOnInit() {
    if (!this.node) {
      let name = this._routeParams.get('name');

      this.unsubscribe = this._appStore.subscribe((state) => {
        this.node = state.consul.node;
        this.isFetchingNode = state.consul.isFetchingNode;
      });

      this._appStore.dispatch(this._consulActions.fetchNode(name));
    }
  }

  goBack() {
    window.history.back();
  }

  private ngOnDestroy() {
    this.unsubscribe();
  }
}
