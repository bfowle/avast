import {Component, View} from 'angular2/core';
import {DockerImage} from '../interfaces/docker-image';

@Component({
    selector: 'image-tree',
    inputs: ['children']
})

@View({
    template: require('./image-tree.directive.html'),
    styles: [
      require('./image-tree.directive.css')
    ],
    directives: [ImageTree]
})

export class ImageTree {
    children: DockerImage[];
}
