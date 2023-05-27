describe('All Mince Pies', () => {
	it('Can visit and see Headers', () => {
		cy.visit('/');
		cy.contains('Current Top Pie');
	});
});
