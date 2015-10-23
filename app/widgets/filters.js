'use strict';

angular.module('dashboard').filter('widgetName', [
  'widgetFactory',
  function(widgetFactory) {
    return function(input) {
      var w = widgetFactory(input.type);
      if(w === null) {
        return 'Unknown';
      }

      if(input.configuration) {
        w.configure(input.configuration);
      }

      return w.name;
    };
  }
]);

angular.module('dashboard').filter('byteSize', [
  function() {
    return function(bytes, precision) {
      if(!precision) {
        precision = 1;
      }

      bytes = parseFloat(bytes);
      if(!bytes) {
        bytes = 0;
      }

      var i, suffixes = ['B', 'kB', 'MB', 'GB', 'PB'];
      for(i = 0; i < suffixes.length; i++) {
        if(bytes < 1024) {
          break;
        }
        bytes /= 1024;
      }

      if(i === suffixes.length) {
        bytes *= 1024;
        i -= 1;
      }

      return bytes.toFixed(precision) + ' ' + suffixes[i];
    };
  }
]);

angular.module('dashboard').filter('bandwidth', [
  'byteSizeFilter',
  function(byteSizeFilter) {
    return function(bytes, time, precision) {
      if(!time) {
        time = 1.0;
      }

      var bytesPerSecond = bytes / time;
      return byteSizeFilter(bytesPerSecond, precision) + '/s';
    };
  }
]);
