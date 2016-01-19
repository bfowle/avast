import {provide} from 'angular2/core';
import {bootstrap, ELEMENT_PROBE_PROVIDERS} from 'angular2/platform/browser';
import {ROUTER_PROVIDERS} from 'angular2/router';
import {HTTP_PROVIDERS} from 'angular2/http';
import {createStore, combineReducers, bindActionCreators, applyMiddleware} from 'redux';
const thunkMiddleware = require('redux-thunk');
import {AppStore} from 'angular2-redux';
import 'rxjs/add/operator/map';

var ReconnectingWebSocket = require('ReconnectingWebSocket');
var socket = new ReconnectingWebSocket(
  `${window.location.origin.replace('http', 'ws')}/ws`, null, { debug: true }
);

socket.onopen = function() {
  console.log('on open');
};

socket.onmessage = function(e) {
  console.log('on message');
  console.log(e);
};

socket.onclose = function() {
  console.log('on close');
};

import {AppComponent} from './app/components/app.component';
import {ConsulActions} from './app/actions/consul.actions';
import {DockerActions} from './app/actions/docker.actions';
import consul from './app/reducers/consul.reducer';
import docker from './app/reducers/docker.reducer';

const loggerMiddleware = store => next => action => {
    console.log('dispatching', action);
    let result = next(action);
    console.log('next state', store.getState());
    return result;
};

let createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  loggerMiddleware
)(createStore);

const initState = {
  consul: {
    //datacenters: ['dc1']
  },
  docker: {
    //containers: [],
    //images: []
  }
};

const appStore = new AppStore(
  createStoreWithMiddleware(combineReducers({
    consul,
    docker
  }), initState)
);

document.addEventListener('DOMContentLoaded', function main() {
  bootstrap(AppComponent, [
    ...(process.env.ENV === 'production' ? [] : ELEMENT_PROBE_PROVIDERS),
    ...HTTP_PROVIDERS,
    ...ROUTER_PROVIDERS,
    provide(AppStore, {useValue: appStore}),
    ConsulActions,
    DockerActions
  ])
  .catch(err => console.error(err));
});
