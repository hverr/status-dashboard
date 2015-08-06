var gulp = require('gulp');
var jshint = require('gulp-jshint');
var recess = require('gulp-recess');

gulp.task('lint', ['jshint', 'recess'], function() {});

gulp.task('jshint', function() {
  gulp.src(["static/app/**/*.js", '!static/app/bower_components/**'])
    .pipe(jshint())
    .pipe(jshint.reporter('default'));
});

gulp.task('recess', function() {
  options = {
    strictPropertyOrder: false,
  };
  gulp.src(["static/app/**/*.less", '!static/app/bower_components/**'])
    .pipe(recess(options))
    .pipe(recess.reporter());
});
