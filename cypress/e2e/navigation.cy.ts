describe('All Mince Pies', () => {
	it('Can visit and see Headers', () => {
		cy.visit('/');
		cy.findAllByTestId('navigation-menu-link').contains('Home');
		cy.findAllByTestId('navigation-menu-link').contains('Home').click();
		cy.url().should('eq', Cypress.config().baseUrl + '/');

		cy.findAllByTestId('navigation-menu-link').contains('All Pies');
		cy.findAllByTestId('navigation-menu-link').contains('All Pies').click();
		cy.url().should('eq', Cypress.config().baseUrl + '/all-mince-pies');

		cy.findAllByTestId('navigation-menu-link').contains('My Rankings');
		cy.findAllByTestId('navigation-menu-link').contains('My Rankings').click();
		cy.url().should('eq', Cypress.config().baseUrl + '/my-rankings');
	});
});
