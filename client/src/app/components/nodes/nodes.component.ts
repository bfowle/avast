import {Component, OnInit, OnDestroy} from 'angular2/core';
import {Router} from 'angular2/router';
import {AppStore} from 'angular2-redux';

import {ConsulActions} from '../../actions/consul.actions';
import {NodeDetailComponent} from './node-detail.component';
import {SwarmNode} from './interfaces/swarm-node';

@Component({
  selector: 'avast-nodes',
  template: require('./nodes.component.html'),
  styles: [
    require('./nodes.component.css')
  ],
  directives: [NodeDetailComponent]
})

/* @TODO: make this a shared directive */

export class NodesComponent implements OnInit {
  public nodes: SwarmNode[];
  public selectedNode: SwarmNode;
  private isFetchingNodes: boolean = false;
  private unsubscribe: Function;

  constructor(private _router: Router,
              private _appStore: AppStore,
              private _consulActions: ConsulActions) {
  }

  ngOnInit() {
    this.unsubscribe = this._appStore.subscribe((state) => {
      this.nodes = state.consul.nodes;
      this.isFetchingNodes = state.consul.isFetchingNodes;
    });

    this._appStore.dispatch(this._consulActions.fetchNodes());
  }

  onSelect(node: SwarmNode) {
    this.selectedNode = node;
  }

  gotoDetail() {
    this._router.navigate(['NodeDetail', {
      name: this.selectedNode.name
    }]);
  }

  private ngOnDestroy() {
    this.unsubscribe();
  }
}
