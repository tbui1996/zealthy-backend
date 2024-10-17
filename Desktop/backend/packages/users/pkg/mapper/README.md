# Usage

This package works in tandem with the `model` package to give an "ORM" feel. Each root mapper, such as ExternalUser and ExternalUserOrganization, exposes methods such as `Find`, `Update`, `Insert`, and `FindAll`, that can be used outside of the package.

## Find

The find method searches the database and returns an object if it exists. If an object is found, the next time the find method is used for that object, it will use the cached version instead of searching the database again.

## Update

The update method is used after having a reference to a domain model. The client business logic is able to use the setter methods on the domain model to change values of the object, and when it's passed to the update method, those changes get interpreted and translated to database operations.

The update method returns a clone of the passed in object that must be used afterward. The cloned object has a different internal state than the original object.

# Design Patterns

## Registry

The entry point to using the mapper package exists in instantiating a `Registry`. The Registry must be used to access the mappers because it knows how to build the mappers, and it's how mappers access each other when they need to.

## Registry

A Registry is a reference to "global" objects. Instead of using global object mechanisms, we pass them down through the registry as dependencies. This allows us to modify the implementation through an interface, and allows the Registry to handle building these global objects.

## Data Mapper

Use a data mapper when we want the database schema and object model to evolve independently and be decoupled from each other. The primary benefit of using the data mapper is that the domain model can ignore the database, both in design and in the build and testing process.

Data Mappers interact with domain models. The business logic uses domain models and then let's mappers interpret what happened to that domain model and translates it to database (or API) operations. At the moment, each domain model has a data mapper, and understands how to translate the operations of a specific domain model to database operations. What this means is, some data mappers may delegate nested domain model operations to other data mappers. For example, on external user update, the ExternalUser data mapper passes the ExternalUserOrganization to the ExternalUserOrganization data mapper update method, rather than handling the update itself. While this may induce performance problems up front, it is simpler to understand through separation of concerns. Additional optimizations or design patterns (such as a unit of work) can be introduced later on.

## Identity Map

The data mappers make use of [Identity Map](https://martinfowler.com/eaaCatalog/identityMap.html) so that client's don't have to concern themselves with caching. This allows for business logic to use `Find` multiple times to find the same object rather than passing down that object through variables.
