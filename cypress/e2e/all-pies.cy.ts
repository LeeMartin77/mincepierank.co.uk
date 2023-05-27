describe('All Mince Pies', () => {
	it('Can visit and see some pies', () => {
		cy.visit('/all-mince-pies');
		cy.contains('Top Pie');
	});
});
