'use strict';

describe('widget filters', function() {
  beforeEach(module('dashboard'));

  describe('widgetName', function() {
    var widgetName;

    beforeEach(inject(function($filter) {
      widgetName = $filter('widgetName');
    }));

    it('should get the name', function() {
      expect(widgetName('load')).to.equal('Load');
    });

    it('should handle unknown widgets', function() {
      expect(widgetName('not existing')).to.equal('Unknown');
    });
  });
});
