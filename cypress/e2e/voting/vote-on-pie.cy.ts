import { signIn } from "../../lib/auth";

describe("Voting on pie", () => {
  it("Can pick a pie, sign in, and return to pie", async () => {
    cy.visit("/all-mince-pies")
    cy.findAllByTestId("pie-summary-card-link-to-pie").first().click();

    cy.findByTestId("pie-page-display-name-header")
    cy.findByTestId("pie-page-picture")
    cy.findByTestId("pie-page-statistics")

    cy.contains('Sign in to Rank')

    const pieUrl = cy.url();

    cy.contains('button', 'Sign in with Development').click()

    signIn();

    cy.url().then(postLogin => {
      pieUrl.should('eq', postLogin);
    })
  })
})
