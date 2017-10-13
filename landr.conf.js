const fs = require('fs')

const getArch = (str) => {
  const [ _, arch ] = str.match(/-([^-]+)\.tar\.gz$/)
  return arch
}

const packagePrettyName = (str) => `Balena for ${getArch(str)}`

const prepAssets = (release) => {
  release.assets = release.assets.map((asset) => {
    return Object.assign({}, asset, {
      prettyName: packagePrettyName(asset.name),
      arch: getArch(asset.name),
      os: 'Linux'
    })
  })

  return release
}

module.exports = {
  theme: 'landr-theme-basic',
  hooks: {
    'post-build': ({ config }) => {
      const data = fs.readFileSync(`${__dirname}/contrib/install.sh`, 'utf-8')
      return fs.writeFileSync(`${config.distDir}/install.sh`, data)
    }
  },
  middleware: (store, action, next) => {
    if (action.type === 'ADD_RELEASE') {
      // intercept all releases and add pretty labels to assets
      action.payload = prepAssets(action.payload)
    }

    return next(action)
  },
  settings: {
    lead: 'A Moby-based container engine for IoT',
    analytics: {
      mixpanelToken: '',  // mixpanelToken
      gosquaredId: '', // gosquared Id
      gaSite: '', // google Analytics site eg balena.io
      gaId: '' // google Analytics ID
    },
    theme: {
      colors: {
        primary: '#00A375'
      }
    },
    callToActionCommand: 'curl -sfL https://balena.io/install.sh | sh',
    features: [
      {
        'title': 'Small footprint',
        'image': 'footprint',
        'description': '3x smaller than Docker CE, packaged as a single binary'
      },
      {
        'title': 'Multi-arch support',
        'image': 'multiple',
        'description': 'Bandwidth-efficient updates with binary diffs, 10-70x smaller than pulling layers in <a href="blog link, hashtagged to the technical comparison title">common scenarios</a>'
      },
      {
        'title': 'True container deltas',
        'image': 'bandwidth',
        'description': 'True container deltas calculate full binary diffs that are 10-50x more bandwidth efficient than the standard layer-based delta updates'
      },
      {
        'title': 'Minimal wear-and-tear',
        'image': 'storage',
        'description': 'Extract layers as they arrive to prevent excessive writing to disk, protecting your storage from eventual corruption'
      },
      {
        'title': 'Failure-resistant pulls',
        'image': 'failure-resistant',
        'description': 'Atomic and durable image pulls defend against partial container pulls in the event of power failure'
      },
      {
        'title': 'Conservative memory use',
        'image': 'undisturbed',
        'description': 'Prevents page cache thrashing during image pull, so your application runs undisturbed in low-memory situations'
      }
    ],
    motivation: [
      'Balena is a new container engine based on Docker’s Moby project, with an emphasis on embedded and IoT use cases, and compatible with Docker containers.</br></br> It supports container deltas for 10-50x more efficient bandwidth usage, has 3x smaller binaries, uses ram and storage more conservatively, and focuses on atomicity and durability of container pulling.'
    ],
  }
}