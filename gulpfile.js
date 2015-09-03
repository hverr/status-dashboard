var gulp = require('gulp');
var jshint = require('gulp-jshint');
var recess = require('gulp-recess');
var trimlines = require('gulp-trimlines');
var uglify = require('gulp-uglify');
var concat = require('gulp-concat');
var rename = require('gulp-rename');
var less = require('gulp-less');
var minify =require('gulp-minify-css');

gulp.task('lint', ['jshint', 'recess'], function() {});

gulp.task('jshint', function() {
  gulp.src(["app/**/*.js", '!app/bower_components/**'])
    .pipe(jshint())
    .pipe(jshint.reporter('default'));
});

gulp.task('recess', function() {
  var options = {
    strictPropertyOrder: false,
  };
  gulp.src(["app/**/*.less", '!app/bower_components/**'])
    .pipe(recess(options))
    .pipe(recess.reporter());
});

gulp.task('trimlines', function() {
  var options = {
    leading: false,
  };
  var source = [
    "app/**/*.less",
    "app/**/*.html",
    "app/**/*.js",
    "!app/bower_components/**",
  ];
  gulp.src(source)
    .pipe(trimlines(options))
    .pipe(gulp.dest("app"));
});

gulp.task('buildjs', function() {
  var source = [
    'app/bower_components/html5-boilerplate/dist/js/vendor/modernizr-2.8.3.min.js',
    'app/bower_components/jquery/dist/jquery.min.js',
    'app/bower_components/bootstrap/dist/js/bootstrap.min.js',
    'app/bower_components/angular/angular.min.js',
    'app/bower_components/angular-route/angular-route.min.js',
    'app/bower_components/javascript-detect-element-resize/jquery.resize.js',
    'app/bower_components/angular-gridster/dist/angular-gridster.min.js',
    'app/app.js',
    'app/filters.js',
    'app/widgets/**/*.js',
  ];

  gulp.src(source)
    .pipe(uglify())
    .pipe(concat('a.js'))
    .pipe(gulp.dest('dist'));
});

gulp.task('buildhtml', function() {
  gulp.src(['app/index_dist.html'])
    .pipe(rename('index.html'))
    .pipe(gulp.dest('dist'));

  var source = [
    'app/**/*.html',

    '!app/index*.html',
    '!app/bower_components/**',
  ];
  gulp.src(source)
    .pipe(gulp.dest('dist'));
});

gulp.task('buildcss', function() {
  gulp.src(['app/app.less'])
    .pipe(less())
    .pipe(concat('a.css'))
    .pipe(minify())
    .pipe(gulp.dest('dist'));

  var source = [
    'app/bower_components/html5-boilerplate/dist/css/normalize.css',
    'app/bower_components/bootstrap/dist/css/bootstrap.min.css',
    'app/bower_components/angular-gridster/dist/angular-gridster.min.css',
  ];

  gulp.src(source)
    .pipe(concat('b.css'))
    .pipe(minify())
    .pipe(gulp.dest('dist'));
});

gulp.task('buildassets', function() {
  gulp.src(['app/bower_components/bootstrap/dist/fonts/**'])
    .pipe(gulp.dest('dist/fonts'));
  gulp.src(['app/favicon.png'])
    .pipe(gulp.dest('dist'));
});

gulp.task('build', ['buildcss', 'buildjs', 'buildhtml', 'buildassets']);
