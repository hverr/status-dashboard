'use strict';

describe('widget directives', function() {
  var $compile;
  var $rootScope;

  beforeEach(module('dashboard'));
  beforeEach(module('testTemplates'));

  beforeEach(inject(function(_$compile_, _$rootScope_) {
    $compile = _$compile_;
    $rootScope = _$rootScope_;
  }));

  describe('widget', function() {
    it('should properly bind', function() {
      var scope = $rootScope.$new();
      scope.widget = {
        client: 'client name',
        name: 'widget name',
      };

      var w = $compile('<div widget></div>')(scope);
      scope.$digest();

      expect(w.html()).to.contain('client name'.toUpperCase());
      expect(w.html()).to.contain('widget name'.toUpperCase());
    });
  });

  describe('widgetDynamicInfo', function() {
    it('should properly replace', function() {
      var scope = $rootScope.$new();
      scope.widget = {
        directive: 'custom-directive',
      };

      $compile('<widget-dynamic-info></widget-dynamic-info>')(scope);
      scope.$digest();
    });
  });
});
