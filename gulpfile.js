const gulp    = require('gulp');
const seq     = require('run-sequence');
const webpack = require('webpack-stream');
const uglify  = require('gulp-uglify');
const del     = require('del');

gulp.task('default', (cb) => seq('clean', 'build', cb));

gulp.task('watch', ['default'], () => {
  gulp.watch('src/main/js/**/*', ['webpack']);
  gulp.watch('static/**/*', ['static-without-reload']);
});

gulp.task('build', ['webpack', 'vendor', 'static']);

gulp.task('webpack', () =>
  gulp.src('src/main/js/main.jsx')
    .pipe(webpack({
      output: { filename: 'main.js' },
      externals: { react: 'React' },
      module: {
        loaders: [
          { test: /\.jsx$/, loader: 'babel-loader?presets[]=es2015&presets[]=react' },
          { test: /\.json$/, loader: 'json-loader' },
          { test: /\.less$/, loader: 'style!css!less?compress' }
        ]
      }
    }))
    .pipe(uglify())
    .pipe(gulp.dest('build/assets')));

gulp.task('vendor', () =>
  gulp.src([
      'node_modules/react/dist/react.min.js',
      'node_modules/bootswatch/cosmo/bootstrap.min.css',
    ]).pipe(gulp.dest('build/assets')));

gulp.task('static', () =>
  gulp.src('static/**/*').pipe(gulp.dest('build/assets')));

// prevent reloading App Engine dev server
gulp.task('static-without-reload', () =>
  gulp.src(['static/**/*', '!static/WEB-INF/appengine-web.xml'])
    .pipe(gulp.dest('build/assets')));

gulp.task('clean', (cb) => del([
  'build/assets/**',
  '!build/assets',
  '!build/assets/WEB-INF',
  // JARs and classes managed by Gradle
  '!build/assets/WEB-INF/lib/**',
  '!build/assets/WEB-INF/classes/**',
  // persistent data
  '!build/assets/WEB-INF/appengine-generated/**',
], cb));
