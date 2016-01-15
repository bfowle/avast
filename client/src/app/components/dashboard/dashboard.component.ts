import {Component, OnInit} from 'angular2/core';
import {Router} from 'angular2/router';
import {AppStore} from 'angular2-redux';

import {ConsulActions} from '../../actions/consul.actions';
import {SwarmNode} from '../nodes/interfaces/swarm-node';

@Component({
  selector: 'avast-dashboard',
  template: require('./dashboard.component.html'),
  styles: [
    require('./dashboard.component.css')
  ]
})

export class DashboardComponent implements OnInit {
  public datacenters: string[] = [];
  public currentDatacenter: string = ''; // @TODO: select first dc on load
  public nodes: SwarmNode[] = [];
  private isFetchingDatacenters: boolean = false;
  private isFetchingNodes: boolean = false;

  constructor (private _router: Router,
               private _appStore: AppStore,
               private _consulActions: ConsulActions) {
  }

  ngOnInit() {
    this._appStore.subscribe((state) => {
      this.datacenters = state.consul.datacenters;
      this.nodes = state.consul.nodes;
      this.isFetchingDatacenters = state.consul.isFetchingDatacenters;
      this.isFetchingNodes = state.consul.isFetchingNodes;
    });

    this._appStore.dispatch(this._consulActions.fetchDatacenters());
  }

  // @TODO: attach this to a change detection (EventEmitter) so that it can be:
  // 1.) loaded by default at the app load
  // 2.) work when someone picks a different dc
  selectDatacenter(dc: string) {
    this.currentDatacenter = dc;
    this._appStore.dispatch(this._consulActions.fetchNodes(this.currentDatacenter));
  }

  gotoNode(name: string) {
    this._router.navigate(['NodeDetail', { name: name }]);
  }
}
