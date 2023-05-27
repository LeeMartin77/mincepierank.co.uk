export function signIn() {
	// This is using default dev credentials
	cy.findByLabelText('Username').type('john');
	cy.findByLabelText('Password').type('doe');

	cy.contains('button', 'Sign in with Development').click();
}
