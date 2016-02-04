'use strict';

const path = require('path');

const descriptor = require('./descriptor');
const Descriptor = descriptor.Descriptor;

const ArigatoConfigError = require('../errors').ArigatoConfigError;
const resolver = require('../util/resolver');

const arigato = exports;

const FILE_NAME = 'arigato.yml';

class Arigato extends Descriptor {
  
  static validate (input) {
    return new Promise((resolve, reject) => {
      input = input || {};
      Descriptor.validate('arigato', input)
        .then(resolve).catch((errors) => {
        if (!Array.isArray(errors) && errors.name !== 'ValidationError') {
          return reject(errors);
        }

        errors = errors.map((err) => {
          return new ArigatoConfigError(err.message, err.code);
        });

        reject(errors);
      });
    });
  }

  static create (base, opts) {
    var filePath = path.join(base, FILE_NAME);
    return super.create(this, 'arigato', filePath, {
      owner: opts.user_id,
      app: opts.app.uuid
    });
  }

  static find (base) {
    return new Promise((resolve, reject) => {
      resolver.parents(base, FILE_NAME).then((files) => {
        // The first one (so index 0) is the closest and therefore most relevant
        if (files.length === 0) {
          return resolve(null);
        }

        var file = files[0];
        Descriptor.read(this, 'arigato', file).then(resolve).catch(reject);
      }).catch(reject);
    });
  }
}

arigato.Arigato = Arigato;
arigato.FILE_NAME = FILE_NAME;