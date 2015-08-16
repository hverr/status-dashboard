'use strict';

angular.module('dashboard').filter('pluralize', [
  function() {
    return function(input, singular, plural) {
      if(input === 1 || input === "1") {
        return singular;
      } else {
        if(plural) {
          return plural;
        } else {
          return singular + 's';
        }
      }
    };
  }
]);
