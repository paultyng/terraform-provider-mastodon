# modify user class to not do dns/mx validation
User.class_eval {
  def validate_email_dns?
    false
  end
}

# create normal users/accounts
application_owner_user = nil
3.times do |i|
  normal_account = Account.create!(username: "acctest#{i}")
  normal_user = User.create!(email: "acctest#{i}@example.com", password: 'password', account: normal_account, agreement: true)
  normal_user.confirm
  normal_account.save!
  normal_user.save!

  application_owner_user ||= normal_user
end

# create admin user/account
admin_account = Account.create!(username: 'acctestadmin')
admin_user = User.create!(email: 'acctest-admin@example.com', password: 'password', account: admin_account, agreement: true, admin: true)
admin_user.confirm
admin_account.save!
admin_user.save!

# site settings

# admin contact
Setting.where(var: :site_contact_username).first_or_initialize(var: :site_contact_username).update(value: admin_account.username)
Setting.where(var: :site_contact_email).first_or_initialize(var: :site_contact_email).update(value: admin_user.email)

# doorkeeper application
# well known secrets for testing
well_known_token = 'terraform-acctest-terraform-acctest-terrafo'

app = Doorkeeper::Application.create!(
  redirect_uri: 'urn:ietf:wg:oauth:2.0:oob',
  name: 'Terraform Acctest',
  owner_type: 'User',
  owner: application_owner_user,
  uid: well_known_token,
  secret: well_known_token,
  scopes: 'read write follow'
)
