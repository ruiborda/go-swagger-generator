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
 go get github.com/ruiborda/go-swagger-generator/v2`}
          </CodeBlock>
        </div>
        <div className="padding-top--md" style={{maxWidth: 650, margin: '0 auto', textAlign: 'left'}}>
          <CodeBlock language="go" className={styles.codePreview}>
            {`// OpenAPI 3.0 documentation for a Go API (using Go Swagger Generator v2)
var _ = swagger.Swagger().Path("/users/{id}"). // Path relative to server base URL
    Get(func(op openapi.Operation) {
        op.Summary("Find user by ID").
            Tag("UserController").
            OperationID("getUserByIdV2").
            PathParameter("id", func(p openapi.Parameter) {
                p.Required(true).
                  Description("User ID").
                  Schema(func(s openapi.Schema) {
                      s.Type("integer").Format("int64")
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                        mt.SchemaFromDTO(UserDto{}) // Assumes UserDto is defined
                    })
            })
    }).
    Doc()

// Gin handler
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   id, // Simplified, parse 'id' to appropriate type
		Name: "John Doe",
	})
}`}
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
{`// Define an endpoint with its OpenAPI 3.0 documentation
var _ = swagger.Swagger().Path("/users/{id}"). // Path relative to server base URL
    Get(func(op openapi.Operation) {
        op.Summary("Find user by ID").
            Tag("UserController").
            OperationID("getUserByIdV2Example").
            PathParameter("id", func(p openapi.Parameter) {
                p.Required(true).
                  Description("User ID").
                  Schema(func(s openapi.Schema) {
                      s.Type("integer").Format("int64")
                  })
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                        mt.SchemaFromDTO(UserDto{}) // Assumes UserDto is defined
                    })
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
