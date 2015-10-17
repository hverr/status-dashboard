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

  describe('byteSize', function() {
    var byteSize;

    beforeEach(inject(function($filter) {
      byteSize = $filter('byteSize');
    }));

    it('should use a precision of at least 1', function() {
      expect(byteSize(5)).to.equal('5.0 B');
      expect(byteSize(12, 0)).to.equal('12.0 B');
      expect(byteSize("12.3", 0)).to.equal('12.3 B');
    });

    it('should handle invalid bytes', function() {
      expect(byteSize('not a number')).to.equal('0.0 B');
      expect(byteSize(0)).to.equal('0.0 B');
    });

    it('should handle B range', function() {
      expect(byteSize(0)).to.equal('0.0 B');
      expect(byteSize(50)).to.equal('50.0 B');
      expect(byteSize(1023)).to.equal('1023.0 B');
    });

    it('should handle kB range', function() {
      expect(byteSize(1024)).to.equal('1.0 kB');
      expect(byteSize(2048)).to.equal('2.0 kB');
    });

    it('should handle MB range', function() {
      expect(byteSize(1*1024*1024)).to.equal('1.0 MB');
      expect(byteSize(2*1024*1024)).to.equal('2.0 MB');
    });

    it('should handle GB range', function() {
      expect(byteSize(1*1024*1024*1024)).to.equal('1.0 GB');
      expect(byteSize(2*1024*1024*1024)).to.equal('2.0 GB');
    });

    it('should handle PB range', function() {
      expect(byteSize(1*1024*1024*1024*1024)).to.equal('1.0 PB');
      expect(byteSize(2*1024*1024*1024*1024)).to.equal('2.0 PB');
    });

    it('should handle >PB range', function() {
      expect(byteSize(2*1024*1024*1024*1024*1024)).to.equal('2048.0 PB');
    });
  });

  describe('bandwidth', function() {
    var bandwidth;

    beforeEach(inject(function($filter) {
      bandwidth = $filter('bandwidth');
    }));

    it('should work correctly', function() {
      expect(bandwidth(2048, 1.0, 2)).to.equal('2.00 kB/s');
      expect(bandwidth(2048, 2.0)).to.equal('1.0 kB/s');
      expect(bandwidth(2*1024*1024)).to.equal('2.0 MB/s');
    });
  });
});
