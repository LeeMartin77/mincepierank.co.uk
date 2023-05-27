describe('loads homepage :: can see expected', () => {
  it('passes', () => {
    cy.visit('/');
    cy.contains('Current Top Pie');
  });
});
