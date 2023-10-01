arar: Assume Role And Run
=========================

arar: Assume Role And Run performs [assume-role](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use.html) and run a command with that credentials.


Context
-------

Combination with AWS_ACCESS_KEY_ID and AWS_PROFILE environment variables doesn't allow assume-role. You must pass AWS_PROFILE without environment variables, e.g. with `--profile`:

```console
# This doesn't work as expected. This just result the identity of the IAM user with AWS_ACCESS_KEY_ID
AWS_ACCESS_KEY_ID=xxx \
AWS_SECRET_ACCESS_KEY=xxx \
AWS_PROFILE=assume-roling-profile \
aws sts get-caller-identity
```

And you cannot specify the role to assume with environment variable. You always need setup `.aws/config` .
(Actually, you should use Web Identity for those situations if you could.)

Those limitations are really hard to use in CI/CD. `arar` provides those mechanisms.


Usage
-----

You can peerform assume-role like this:

```console
AWS_ACCESS_KEY_ID=xxx \
AWS_SECRET_ACCESS_KEY=xxx \
AWS_REGION=us-east-1 \
AWS_ROLE_ARN=arn:aws:iam::xxxxxxxxxx:role/assumerole-role \
AWS_ROLE_SESSION_NAME=assumerole-user \
arar -- aws sts get-caller-identity
```

Also supports command line options (`aws sts assume-role` compatible):

```console
AWS_ACCESS_KEY_ID=xxx \
AWS_SECRET_ACCESS_KEY=xxx \
arar \
  --region=ap-northeast-1 \
  --role-arn=arn:aws:iam::xxxxxxxxxxx:role/assumerole-role \
  --role-session-name=assumerole-user
  -- \
  aws sts get-caller-identity
```

MFA support
-----------

Supports token codes with MFA devices:

```console
arar \
  --serial-number=arn:aws:iam::xxxxxxxxxxx:mfa/assumeroling-user
  -- \
  aws sts get-caller-identity
```

Automatic session name with IAM user name
-----------------------------------------

`-u/--username-session` option sets the session name with IAM user name.
This is useful when the role is configured with `"sts:RoleSessionName": "${aws:username}"`. ([AWS Security Blog post](https://aws.amazon.com/jp/blogs/security/easily-control-naming-individual-iam-role-sessions/)).

```console
arar -u -- aws sts get-caller-identity
```
