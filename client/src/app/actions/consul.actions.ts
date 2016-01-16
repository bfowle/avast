import {Injectable} from 'angular2/core';
import {Http, Response} from 'angular2/http';
import {Actions} from 'angular2-redux';
import 'rxjs/add/operator/map';

import {API_VERSION} from '../constants';
import {SwarmNode} from '../components/nodes/interfaces/swarm-node';

export const REQUEST_DATACENTERS = 'REQUEST_DATACENTERS';
export const RECEIVE_DATACENTERS = 'RECEIVE_DATACENTERS';
export const REQUEST_NODES = 'REQUEST_NODES';
export const RECEIVE_NODES = 'RECEIVE_NODES';
export const REQUEST_NODE = 'REQUEST_NODE';
export const RECEIVE_NODE = 'RECEIVE_NODE';

@Injectable()
export class ConsulActions extends Actions {

  constructor(private _http: Http) {
    super();
  }

  fetchDatacenters() {
    return (dispatch) => {
      dispatch(this.requestDatacenters());

      this._http.get(`/api/${API_VERSION}/consul/datacenters`)
        .map((res: Response) => res.json())
        .map(res => dispatch(this.receiveDatacenters(res)))
        .subscribe();
    };
  }

  fetchNodes(dc: string = '') {
    return (dispatch) => {
      dispatch(this.requestNodes(dc));

      this._http.get(`/api/${API_VERSION}/consul/nodes${dc ? '/' + dc : ''}`)
        .map((res: Response) => res.json())
        .map(res => dispatch(this.receiveNodes(dc, res)))
        .subscribe();
    };
  }

  fetchNode(name: string) {
    return (dispatch) => {
      dispatch(this.requestNode(name));

      this._http.get(`/api/${API_VERSION}/consul/node/${name}`)
        .map((res: Response) => res.json())
        .map(res => dispatch(this.receiveNode(name, res)))
        .subscribe();
    };
  }

  requestDatacenters() {
    return {type: REQUEST_DATACENTERS};
  }

  requestNodes(dc: string) {
    return {type: REQUEST_NODES};
  }

  requestNode(name: string) {
    return {type: REQUEST_NODE};
  }

  receiveDatacenters(datacenters: string[]) {
    return {
      type: RECEIVE_DATACENTERS,
      datacenters
    };
  }

  receiveNodes(dc: string, nodes: SwarmNode[]) {
    return {
      type: RECEIVE_NODES,
      dc,
      nodes
    };
  }

  receiveNode(name: string, node: SwarmNode) {
    return {
      type: RECEIVE_NODE,
      name,
      node
    };
  }

}
