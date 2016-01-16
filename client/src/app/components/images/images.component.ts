import {Component, OnInit, OnDestroy} from 'angular2/core';
import {Router} from 'angular2/router';
import {AppStore} from 'angular2-redux';

import {DockerActions} from '../../actions/docker.actions';
import {DockerImage} from './interfaces/docker-image';

@Component({
  selector: 'avast-images',
  template: require('./images.component.html'),
  styles: [
    require('./images.component.css')
  ]
})

export class ImagesComponent implements OnInit {
  public rootImage: DockerImage;
  public selectedImage: DockerImage;
  private isFetchingImages: boolean = false;
  private unsubscribe: Function;

  constructor(private _router: Router,
              private _appStore: AppStore,
              private _dockerActions: DockerActions) {
  }

  ngOnInit() {
    this.unsubscribe = this._appStore.subscribe((state) => {
      this.rootImage = state.docker.rootImage;
      this.isFetchingImages = state.docker.isFetchingImages;
    });

    this._appStore.dispatch(this._dockerActions.fetchImages());
  }

  onSelect(image: DockerImage) {
    this.selectedImage = image;
  }

  private ngOnDestroy() {
    this.unsubscribe();
  }
}
