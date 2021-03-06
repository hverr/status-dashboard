module.exports = function(config) {
  config.set({
    basePath: '',
    frameworks: ['mocha', 'chai'],
    files: [
      'app/bower_components/jquery/dist/jquery.min.js',
      'app/bower_components/bootstrap/dist/js/bootstrap.min.js',
      'app/bower_components/angular/angular.js',
      'app/bower_components/angular-route/angular-route.js',
      'app/bower_components/angular-cookies/angular-cookies.js',
      'app/bower_components/javascript-detect-element-resize/jquery.resize.js',
      'app/bower_components/angular-gridster/dist/angular-gridster.min.js',

      'node_modules/angular-mocks/angular-mocks.js',

      'app/app.js',
      'app/filters.js',
      'app/widgets/**/**.js',
      'app/**/*.html',

      'app/tests/**/*.js',
    ],
    preprocessors: {
      'app/app.js' : ['coverage'],
      'app/filters.js' : ['coverage'],
      'app/widgets/**/**.js' : ['coverage'],
      'app/**/*.html' : ['ng-html2js'],
    },
    ngHtml2JsPreprocessor: {
      stripPrefix: 'app/',
      moduleName: 'testTemplates',
    },
    reporters: ['progress', 'coverage'],
    port: 9876,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: false,
    browsers: ['PhantomJS'],
    singleRun: true,
  });
};
