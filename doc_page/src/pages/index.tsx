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
          <CodeBlock language="go" className={styles.codePreview}>
            {`// Swagger documentation for a Go API
var _ = swagger.Swagger().Path("/users/{id}").
    Get(func(op openapi.Operation) {
        op.Summary("Find user by ID").
            Tag("UserController").
            Produce(mime.ApplicationJSON).
            PathParameter("id", func(p openapi.Parameter) {
                p.Required(true).Type("integer")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    SchemaFromDTO(UserDto{})
            })
    }).
    Doc()
// Gin handler
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"id":   id,
		"name": "John Doe",
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
              Generate Swagger Documentation with Go
            </Heading>
            <p className="text--center padding-vert--md">
              Go Swagger Generator lets you document your Go APIs with Swagger/OpenAPI in minutes
            </p>
            
            <div className="padding-top--md" style={{maxWidth: 700, margin: '0 auto'}}>
              <CodeBlock language="go" className={styles.codePreview}>
{`// Define an endpoint with its Swagger documentation
var _ = swagger.Swagger().Path("/users/{id}").
    Get(func(op openapi.Operation) {
        op.Summary("Find user by ID").
            Tag("UserController").
            Produce(mime.ApplicationJSON).
            PathParameter("id", func(p openapi.Parameter) {
                p.Required(true).Type("integer")
            }).
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("successful operation").
                    SchemaFromDTO(UserDto{})
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
      title={`${siteConfig.title} - Swagger Documentation Generator for Go`}
      description="Go Swagger Generator - A library for generating Swagger/OpenAPI documentation for Go APIs">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <QuickStartSection />
      </main>
    </Layout>
  );
}
