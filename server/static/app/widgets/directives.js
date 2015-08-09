'use strict';

angular.module('dashboard').directive('widget', [
  function() {
    return {
      templateUrl : 'widgets/widget.html',
    };
  }
]);
