const gulp = require('gulp');
const sass = require('gulp-sass');
const watch = require('gulp-watch');
const batch = require('gulp-batch');
const plumber = require('gulp-plumber');
const jetpack = require('fs-jetpack');
const bundle = require('./bundle');
const utils = require('./utils');

const projectDir = jetpack;
const srcDir = jetpack.cwd('./src');
const destDir = jetpack.cwd('./app');

gulp.task('bundle', () => {
  return Promise.all([
    bundle(srcDir.path('background.js'), destDir.path('background.js')),
    bundle(srcDir.path('app.js'), destDir.path('app.js')),
  ]);
});

gulp.task('sass', () => {
  return gulp.src(srcDir.path('stylesheets/main.scss'))
  .pipe(plumber())
  .pipe(sass().on('error', sass.logError))
  .pipe(gulp.dest(destDir.path('stylesheets')));
});

gulp.task('environment', () => {
  const configFile = `config/env_${utils.getEnvName()}.json`;
  projectDir.copy(configFile, destDir.path('env.json'), { overwrite: true });
});

gulp.task('watch', () => {
  const beepOnError = (done) => {
    return (err) => {
      if (err) {
        utils.beepSound();
      }
      done(err);
    };
  };

  watch('src/**/*.js', batch((events, done) => {
    gulp.start('bundle', beepOnError(done));
  }));
  watch('src/**/*.scss', batch((events, done) => {
    gulp.start('sass', beepOnError(done));
  }));
});

gulp.task('build', ['bundle', 'sass', 'environment']);
