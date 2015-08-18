'use strict';

describe('meminfo widget', function() {
  var MeminfoWidget;

  beforeEach(module('dashboard'));
  beforeEach(module('testTemplates'));

  beforeEach(inject(function(_MeminfoWidget_) {
    MeminfoWidget = _MeminfoWidget_;
  }));

  describe('MeminfoWidget', function() {
    it('should construct properly', function() {
      var w = new MeminfoWidget();

      expect(w.directive).to.equal('meminfo-widget');
      expect(w.name).to.equal('Memory');

      expect(w.data).to.have.property('total');
      expect(w.data).to.have.property('used');
    });
  });
});
