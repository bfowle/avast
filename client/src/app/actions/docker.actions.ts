import {Injectable} from 'angular2/core';
import {Http, Response} from 'angular2/http';
import {Actions} from 'angular2-redux';
import 'rxjs/add/operator/map';

import {API_VERSION} from '../constants';
import {DockerContainer} from '../components/containers/interfaces/docker-container';
import {DockerImage} from '../components/images/interfaces/docker-image';

export const REQUEST_CONTAINERS = 'REQUEST_CONTAINERS';
export const RECEIVE_CONTAINERS = 'RECEIVE_CONTAINERS';
export const REQUEST_IMAGES = 'REQUEST_IMAGES';
export const RECEIVE_IMAGES = 'RECEIVE_IMAGES';

@Injectable()
export class DockerActions extends Actions {

  constructor(private _http: Http) {
    super();
  }

  fetchContainers() {
    return (dispatch) => {
      dispatch(this.requestContainers());

      this._http.get(`/api/${API_VERSION}/docker/containers`)
        .map((res: Response) => res.json())
        .map(res => dispatch(this.receiveContainers(res)))
        .subscribe();
    };
  }

  fetchImages() {
    return (dispatch) => {
      dispatch(this.requestImages());

      this._http.get(`/api/${API_VERSION}/docker/images`)
        .map((res: Response) => res.json())
        .map(res => dispatch(this.receiveImages(res)))
        .subscribe();
    };
  }

  requestContainers() {
    return {type: REQUEST_CONTAINERS};
  }

  requestImages() {
    return {type: REQUEST_IMAGES};
  }

  receiveContainers(containers: DockerContainer[]) {
    return {
      type: RECEIVE_CONTAINERS,
      containers
    };
  }

  receiveImages(rootImage: DockerImage) {
    return {
      type: RECEIVE_IMAGES,
      rootImage
    };
  }

}
