import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';
import CodeBlock from '@theme/CodeBlock';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <div className="padding-top--md" style={{maxWidth: 650, margin: '0 auto', textAlign: 'left'}}>
          <CodeBlock language="bash" className={styles.codePreview}>
            {`# Instala Go Swagger Generator v2
 go get -u github.com/ruiborda/go-swagger-generator/v2`}
          </CodeBlock>
        </div>
        <div className="padding-top--md" style={{maxWidth: 650, margin: '0 auto', textAlign: 'left'}}>
          <CodeBlock language="go" className={styles.codePreview}>
            {`// UserDto represents the data transfer object for a user.
type UserDto struct {
    ID   int    ${'`'}json:"id" yaml:"id"${'`'}
    Name string ${'`'}json:"name" yaml:"name"${'`'}
}

// GetUserByIdHandler is the Gin handler for retrieving a user.
func GetUserByIdHandler(c *gin.Context) {
    idStr := c.Param("id")
    c.JSON(http.StatusOK, UserDto{
        ID:   1, // Example ID
        Name: "John Doe (User " + idStr + ")",
    })
}

// --- In your ConfigureOpenAPI function (see Quick Start for full example) ---
// doc := swagger.Swagger() // Get OpenAPI builder instance
// _, _ = doc.SchemaFromDTO(&UserDto{}) // Register DTO

// Document the /users/{id} GET endpoint (relative to server URL)
var _ = doc.Path("/users/{id}").
    Get(func(op openapi.Operation) {
        op.Summary("Find user by ID").
            Tag("User Management").
            OperationID("getUserById").
            PathParameter("id", func(p openapi.Parameter) {
                p.Description("ID of the user to retrieve").
                    Required(true).
                    Schema(func(s openapi.Schema) {
                        s.Type("integer").Format("int64")
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - user details returned").
                    Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                        mt.SchemaFromDTO(&UserDto{})
                    })
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()`}
          </CodeBlock>
        </div>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/quick-start">
            Quick Start 🚀
          </Link>
        </div>
      </div>
    </header>
  );
}

function QuickStartSection() {
  return (
    <section className={clsx('padding-vert--xl', styles.quickStart)}>
      <div className="container">
        <div className="row">
          <div className="col col--12">
            <Heading as="h2" className="text--center">
              Generate OpenAPI 3.0 Documentation with Go
            </Heading>
            <p className="text--center padding-vert--md">
              Go Swagger Generator v2 lets you document your Go APIs with OpenAPI 3.0 in minutes
            </p>
            
            <div className="padding-top--md" style={{maxWidth: 700, margin: '0 auto'}}>
              <CodeBlock language="go" className={styles.codePreview}>
{`// Document the /users/{id} GET endpoint.
// The path "/users/{id}" is relative to the server URL (e.g., http://localhost:8080/v1/users/{id}).
// Assumes 'doc' is your swagger.Swagger() instance and UserDto is defined.
var _ = doc.Path("/users/{id}").
    Get(func(op openapi.Operation) {
        op.Summary("Find user by ID").
            Tag("User Management").
            OperationID("getUserById").
            PathParameter("id", func(p openapi.Parameter) {
                p.Description("ID of the user to retrieve").
                    Required(true).
                    Schema(func(s openapi.Schema) {
                        s.Type("integer").Format("int64")
                    })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - user details returned").
                    Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                        mt.SchemaFromDTO(&UserDto{}) // Assumes UserDto is defined
                    })
            }).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
    }).
    Doc()`}
              </CodeBlock>
            </div>
            
            <div className="padding-top--xl">
              <p className="text--center">
                <Link
                  className="button button--primary button--lg"
                  to="/docs/quick-start">
                  View Quick Start Guide
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - OpenAPI 3.0 Documentation Generator for Go`}
      description="Go Swagger Generator v2 - A library for generating OpenAPI 3.0 documentation for Go APIs">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <QuickStartSection />
      </main>
    </Layout>
  );
}
