'use strict';

describe('widget filters', function() {
  beforeEach(module('dashboard'));

  describe('widgetName', function() {
    var widgetName;

    beforeEach(inject(function($filter) {
      widgetName = $filter('widgetName');
    }));

    it('should get the name', function() {
      expect(widgetName({type:'load'})).to.equal('Load');
    });

    it('should handle configuration', function() {
      var c;

      c = {
        type : 'network',
      };
      expect(widgetName(c)).to.equal('Network');

      c = {
        type : 'network',
        configuration : {
          'interface' : 'eth0',
        },
      };
      expect(widgetName(c)).to.equal('Network (eth0)');
    });

    it('should handle unknown widgets', function() {
      expect(widgetName('not existing')).to.equal('Unknown');
    });
  });
});
