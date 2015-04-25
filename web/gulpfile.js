var gulp = require('gulp');
var addsrc = require('gulp-add-src');
var rename = require("gulp-rename");
var del = require('del');
var sass = require('gulp-sass');
var sourcemaps = require('gulp-sourcemaps');
var plumber = require('gulp-plumber');
var bower = require('gulp-bower');
var concat = require('gulp-concat');
var jshint = require('gulp-jshint');
var jshintreporter = require('jshint-stylish');
var uglify = require('gulp-uglify');
var ngAnnotate = require('gulp-ng-annotate');
var html2js = require('gulp-html2js');
var runSequence = require('run-sequence');

var paths = {
  src: {
    js: [
    'src/utils/utils.js',
    'src/home/home.js',
    'src/job/job.js',
    'src/project/project.js',
    'src/home/home.js',
    'src/**/*.js',
    '!src/dev.js'],
    scss: 'src/**/*.scss',
    html: ['src/**/*.html', '!index*.html'],
    htmlIndex: 'src/index-build.html',
    images: 'src/images/*'
  },
  deps: {
    js: ['vendor/jquery/dist/jquery.js', 'vendor/angular/angular.js', 'vendor/angular-route/angular-route.js', 'vendor/moment/moment.js', 'vendor/lodash/dist/lodash.js'],
    css: ['vendor/pure/pure.css', 'vendor/pure/grids-responsive.css', 'vendor/fontawesome/css/font-awesome.css'],
    fonts: 'vendor/fontawesome/fonts/*'
  },
  dest: {
    root: 'build/web/',
    js: 'build/web/js/',
    css: 'build/web/css/',
    fonts: 'build/web/fonts/',
    images: 'build/web/images'
  }
};

gulp.task('clean:build', function (cb) {
  del([paths.dest.root+'**'], cb);
});

gulp.task('clean:bower', function (cb) {
  del(['vendor/**'], cb);
});

gulp.task('bower', ['clean:bower'], function() {
  return bower()
    .pipe(gulp.dest('vendor/'))
});

gulp.task('js:hint', function () {
    gulp.src(paths.src.js)
      .pipe(jshint({
        'lookup': false,
        '-W097': false,
        'predef': ['$', '_', 'angular', 'moment', 'console']
      }))
      .pipe(jshint.reporter(jshintreporter))
      .pipe(jshint.reporter('fail'));
});


gulp.task('js:app', ['js:hint', 'clean:build'], function () {
    gulp.src(paths.src.html)
        .pipe(html2js({
          outputModuleName: 'bzk.templates',
          useStrict: true,
          base: 'src/'
        }))
        .pipe(concat('bzk-templates.js'))
        .pipe(addsrc(paths.src.js))
        .pipe(ngAnnotate())
        .pipe(uglify())
        .pipe(concat('bzk.js'))
        // .pipe(size())
        .pipe(gulp.dest(paths.dest.js));
});

gulp.task('js:vendor', ['bower'], function () {
    gulp.src(paths.deps.js)
        .pipe(uglify())
        .pipe(concat('vendor.js'))
        // .pipe(size())
        .pipe(gulp.dest(paths.dest.js));
});

// Compile sass files
gulp.task('css:sass', function () {
    gulp.src(paths.src.scss)
    	.pipe(plumber())
    	.pipe(sourcemaps.init())
        .pipe(sass())
        .pipe(sourcemaps.write())
        .pipe(plumber.stop())
        .pipe(gulp.dest(paths.dest.css));
});

gulp.task('css:app', ['clean:build'], function(){
  gulp.src(paths.src.scss)
      .pipe(sass())
      .pipe(concat('bzk.css'))
      .pipe(gulp.dest(paths.dest.css));
});

gulp.task('css:vendor', ['bower'], function () {
    gulp.src(paths.deps.css)
        .pipe(concat('vendor.css'))
        .pipe(gulp.dest(paths.dest.css));
});

gulp.task('assets:app', ['clean:build'], function () {
    gulp.src(paths.src.images)
        .pipe(gulp.dest(paths.dest.images));

    gulp.src(paths.src.htmlIndex)
      .pipe(rename('index.html'))
      .pipe(gulp.dest(paths.dest.root));

});

gulp.task('assets:vendor', ['bower', 'clean:build'], function () {
    gulp.src(paths.deps.fonts)
        .pipe(gulp.dest(paths.dest.fonts));
});


gulp.task('build', [
  'js:app', 'js:vendor',
  'css:app', 'css:vendor',
  'assets:app', 'assets:vendor'
  ]);

// The default task (called when you run `gulp`)
gulp.task('default', ['css:sass'], function() {
  // Watch files and run tasks if they change

  gulp.watch([paths.src.scss], ['css:sass']);

});
