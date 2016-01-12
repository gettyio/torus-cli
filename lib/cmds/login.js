'use strict';

const inquirer = require('inquirer');

const log = require('../util/log').get('cmds/login');
const CommandInterface = require('../command').CommandInterface;
const validation = require('../util/validation');
const users = require('../users');

const QUESTIONS = [
  validation.email,
  validation.password
];

class Login extends CommandInterface {

  execute () {
    return new Promise((resolve, reject) => {
      users.loggedIn().then((loggedIn) => {
        if (!loggedIn) {
          return Promise.resolve();
        }

        log.warn('You are currently logged in, logging you out!');
        return users.logout();
      }).then(() => {
        inquirer.prompt(QUESTIONS, (answers) => {
          users.login(answers).then((results) => {
            /*jshint unused:false*/
            log.print('You have logged in!');
            resolve(true);
          }).catch(reject);
        });
      }).catch(reject);
    });
  }
}

module.exports = {
  key: 'login',
  synopsis: 'logins into a user accout',
  usage: 'arigato login',
  Command: Login,
  example: `\tlocalhost$ arigato login`
};