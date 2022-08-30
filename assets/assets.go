package assets

import "embed"

//go:generate npm --prefix ./js --quiet --no-progress --loglevel=error --no-fund --no-audit install
//go:generate esbuild --bundle --minify --sourcemap --outfile=static/js/bundle.js --define:process.env.NODE_ENV="production" js/src/index.ts
//go:generate js/node_modules/sass/sass.js --load-path=js/node_modules --no-source-map stylesheets/application.scss static/styles/application.css

//go:embed static
var Assets embed.FS
