# Piescraper

## The problem

Most supermarkets and pie-sellers have some form of bot protection on their website - totally valid.

If this was a big commercial project, we'd pay for access, but for a hobby project like this, it makes sense to just copy down the data to a CSV and upload it - which is at least progress against the edit-each-pie-individually model of the past.

Problem is this is massively boring gruntwork, so the goal is to simplify it a bit.

# Extract

We'll have a "webscraping" script that does two things - dumps out the html and images from pages as we're browsing, but does so with names that help an LLM figure out what website it's from.

# Transform

We then use some simple models - florence to identify which pictures are actually of food, and whatever free model we fancy to just sift through the data, combine it with some template data for labels and brands.

# Load

Drop the data in the required structured zip into the website, and we've got a year of mince pies ready to go.


## Considerations

- We want a script that doesn't tie us to a one-and-done upload - the prompts and structure should focus on the idea we might be coming back later to add more (seasonal, after all)