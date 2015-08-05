var gulp = require('gulp');
var jshint = require('gulp-jshint');

gulp.task('hint', function() {
  gulp.src(["static/app/**/*.js", '!static/app/bower_components/**'])
    .pipe(jshint())
    .pipe(jshint.reporter('default'));
});
