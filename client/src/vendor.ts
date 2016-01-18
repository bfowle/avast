// polyfills
import 'es6-shim';
// these modules are what are in 'angular2/bundles/angular2-polyfills' so don't use that here
import 'es6-promise';
import 'zone.js/lib/browser/zone-microtask';

if (process.env.ENV !== 'production') {
  // reflect polyfill
  require('es7-reflect-metadata/dist/browser');
  Error['stackTraceLimit'] = Infinity;
  Zone['longStackTraceZone'] = require('zone.js/lib/zones/long-stack-trace.js');
}

if (process.env.ENV === 'production') {
  // reflect with es7-reflect-metadata/reflect-metadata is added
  // by webpack.prod.config ProvidePlugin
  let ngCore = require('angular2/core');
  ngCore.enableProdMode();
}

// angular2
import 'angular2/platform/browser';
import 'angular2/platform/common_dom';
import 'angular2/core';
import 'angular2/router';
import 'angular2/http';

// rxjs
import 'rxjs';
