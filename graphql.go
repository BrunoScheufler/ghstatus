package main

// GitHub types

// Queries & Mutations

var statusFragment = `
fragment StatusFragment on UserStatus {
	createdAt
	updatedAt
	expiresAt
	message
	emoji
	indicatesLimitedAvailability
	organization {
		name
	}
}
`

var retrievalQuery = `
query StatusRetrievalQuery {
  viewer {
    status {
      ...StatusFragment
    }
  }
}
` + statusFragment

var updateMutation = `
mutation UpdateUserStatusMutation ($newStatus: ChangeUserStatusInput!) {
  changeUserStatus(input: $newStatus) {
    status {
      ...StatusFragment
    }
  }
}
` + statusFragment
