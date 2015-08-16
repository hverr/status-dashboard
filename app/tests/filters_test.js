'use strict';

describe('general filter', function() {

  beforeEach(module('dashboard'));

  describe('pluralize', function() {
    var pluralize;

    beforeEach(inject(function($filter) {
      pluralize = $filter('pluralize');
    }));

    it('should not pluralize 1', function() {
      expect(pluralize(1, "s", "p")).to.equal("s");
      expect(pluralize("1", "s", "p")).to.equal("s");
    });

    it('should choose plural argument', function() {
      expect(pluralize(0, "s", "p")).to.equal("p");
      expect(pluralize(5, "s", "p")).to.equal("p");
    });

    it('should automatically pluralize', function() {
      expect(pluralize(4, "car")).to.equal("cars");
      expect(pluralize(0, "dog")).to.equal("dogs");
    });
  });
});
