'use strict';

describe('widget services', function() {
  beforeEach(module('dashboard'));

  describe('Widget', function() {
    it('should construct properly', inject(function(Widget) {
      var w = new Widget('my-directive', 'Widget Name');

      expect(w.directive).to.equal('my-directive');
      expect(w.name).to.equal('Widget Name');

      expect(w.height).to.equal(1);
      expect(w.width).to.equal(1);
      expect(w.row).to.equal(0);
      expect(w.col).to.equal(0);

      expect(w.type).to.equal(null);
      expect(w.clientIdentifier).to.equal(null);
      expect(w.client).to.equal(null);

      expect(w.data).to.equal(null);
    }));
  });
});
